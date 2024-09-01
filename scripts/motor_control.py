#!/usr/bin/python3
'''
Calcalutes the state of motor control based on the state of the motor and
points control
'''

import json
import sys

THREASHOLD = 5


def main() -> dict[str, str]:
    '''
    Main subroutine. Returns signal output state
    '''

    result: dict[str, str] = {}
    topics: dict[str, str] = json.load(sys.stdin)

    topic_objects = {
        topic: json.loads(value) for topic, value in topics.items()
    }

    if not topic_objects.get('layout/i2c-agent/input/41'):
        return {}

    if len(topic_objects['layout/i2c-agent/input/41']) < 2:
        return {}

    motor_input: list[int] = topic_objects['layout/i2c-agent/input/41']

    if not topic_objects.get('layout/i2c-agent/output/40'):
        return {}

    if len(topic_objects['layout/i2c-agent/output/40']) < 2:
        return {}

    points: list[int] = topic_objects['layout/i2c-agent/output/40']

    if not topic_objects.get('layout/i2c-agent/output/41'):
        return {}

    if len(topic_objects['layout/i2c-agent/output/41']) < 3:
        return {}

    motor_output: list[int] = topic_objects['layout/i2c-agent/output/41']

    presence_sensor: tuple[dict[bool]] = [
        topic_objects.get('layout/agent-60/state'),
        topic_objects.get('layout/agent-61/state')
    ]

    if motor_input[1] > THREASHOLD:
        if (motor_output[0] & 0x4) and (points[0] & 0x2):
            motor_output[0] = 5
            motor_output[1] = motor_output[2]
        elif (motor_output[0] & 0x8) and (points[0] & 0x8):
            motor_output[0] = 9
            motor_output[1] = motor_output[2]
    elif (motor_input[0] > THREASHOLD) and (motor_output[0] & 0x2):
        if (points[0] & 0x2) and (presence_sensor[0]["detector0"] is False):
            motor_output[0] = 10
            motor_output[2] = motor_output[1]
        elif (points[0] & 0x8) and (presence_sensor[1]["detector1"] is False):
            motor_output[0] = 6
            motor_output[2] = motor_output[1]

    result_str = json.dumps(motor_output)
    if result_str != topics['layout/i2c-agent/output/41']:
        result['layout/i2c-agent/output/41'] = result_str

    return result


if __name__ == '__main__':
    print(json.dumps(main()))
