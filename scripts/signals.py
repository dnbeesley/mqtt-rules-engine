#!/usr/bin/python3
'''
Calcalutes the state of signals based on the state of the motor and points
control
'''

import json
import sys

OUTPUTS = {
    "layout/agent-60/signal": {
        0x1: {
            0x1: "GRRR",
            0x2: "RRG0"
        },
        0x2: {
            0x5: "RGRR",
            0x6: "RYY0",
            0xA: "RRYG"
        }
    },
    "layout/agent-61/signal": {
        0x4: {
            0x1: "RGRR",
            0x2: "RR0G"
        },
        0x8: {
            0x6: "RRGY",
            0x9: "GRRR",
            0xA: "YR0Y"
        }
    },
    "layout/agent-62/signal": {
        0x14: {
            0x8: "Y0RR"
        },
        0x18: {
            0x8: "G0RR"
        },
        0x60: {
            0x4: "RRGR",
            0x8: "YGRR"
        },
        0xA0: {
            0x4: "RRRG",
            0x8: "YGRR"
        }
    }
}


def main() -> dict[str, str]:
    '''
    Main subroutine. Returns signal output state
    '''

    result: dict[str, str] = {}
    topics: dict[str, str] = json.load(sys.stdin)

    topic_objects = {
        topic: json.loads(value) for topic, value in topics.items()
    }

    if not topic_objects.get('layout/i2c-agent/output/40'):
        return {}

    if len(topic_objects['layout/i2c-agent/output/40']) < 2:
        return {}

    points: int = topic_objects['layout/i2c-agent/output/40'][0]
    points += 0x100 * topic_objects['layout/i2c-agent/output/40'][1]

    if not topic_objects.get('layout/i2c-agent/output/41'):
        return {}

    if len(topic_objects['layout/i2c-agent/output/41']) < 3:
        return {}

    motor: int = topic_objects['layout/i2c-agent/output/41'][0]

    for topic, output in OUTPUTS.items():
        for points_selector, motor_selectors in output.items():
            if points_selector & points == points_selector:
                for motor_selector, value in motor_selectors.items():
                    if motor_selector & motor == motor_selector:
                        result[topic] = value
                        break

        if topic not in result:
            result[topic] = 'RRRR'

    if topic_objects['layout/i2c-agent/output/41'][1] > 0:
        result['layout/agent-63/signal'] = 'ON'
    else:
        result['layout/agent-63/signal'] = 'OFF'

    return result


if __name__ == '__main__':
    print(json.dumps(main()))
