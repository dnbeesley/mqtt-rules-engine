#!/usr/bin/python3
'''
Calcalutes the state of motor control based on the state of the motor and
points control
'''

import json
import sys
import traceback

THREASHOLD = 5


def main() -> dict[str, str]:
    '''
    Main subroutine. Returns signal output state
    '''

    try:
        result: dict[str, str] = {}
        topics: dict[str, str] = json.load(sys.stdin)

        topic_objects = {
            topic: json.loads(value) for topic, value in topics.items()
        }

        motor_input: list[int] = topic_objects['layout/i2c-agent/input/41']
        if len(motor_input) < 2:
            raise Exception('To few elements in motor input')

        points: list[int] = topic_objects['layout/i2c-agent/output/40']
        if len(points) < 2:
            raise Exception('To few elements in points input')

        motor_output: list[int] = topic_objects['layout/i2c-agent/output/41']
        if len(motor_output) < 3:
            raise Exception('To few elements in points input')

        presence_sensor: tuple[dict[bool]] = [
            topic_objects.get('layout/agent-60/state'),
            topic_objects.get('layout/agent-61/state'),
            topic_objects.get('layout/agent-62/state')
        ]

        if motor_input[1] > THREASHOLD:
            if motor_output[0] & 0x4:
                if points[0] & 0x2:
                    motor_output[0] = 5
                    motor_output[1] = motor_output[2]
                elif presence_sensor[2]["detector0"] is False:
                    motor_output[2] = motor_output[1] = 0
            elif motor_output[0] & 0x8:
                if points[0] & 0x8:
                    motor_output[0] = 9
                    motor_output[1] = motor_output[2]
                elif presence_sensor[2]["detector1"] is False:
                    motor_output[2] = motor_output[1] = 0
        elif motor_input[0] > THREASHOLD:
            if motor_output[0] & 0x1:
                if (points[0] & 0x2) and (presence_sensor[0]["detector1"] is False):
                    motor_output[2] = motor_output[1] = 0
                elif (points[0] & 0x8) and (presence_sensor[1]["detector0"] is False):
                    motor_output[2] = motor_output[1] = 0
            elif motor_output[0] & 0x2:
                if (points[0] & 0x2) and (presence_sensor[0]["detector0"] is False):
                    motor_output[0] = 10
                    motor_output[2] = motor_output[1]
                elif (points[0] & 0x8) and (presence_sensor[1]["detector1"] is False):
                    motor_output[0] = 6
                    motor_output[2] = motor_output[1]

        result_str = json.dumps(motor_output)
        if result_str != topics['layout/i2c-agent/output/41']:
            result['layout/i2c-agent/output/41'] = result_str
    except Exception:
        traceback.print_exc()
        print(str(topics), file=sys.stderr)
        result = {}

    return result


if __name__ == '__main__':
    print(json.dumps(main()))
