import typing


def get_named_args_validator(expected_args: typing.List[str]):
    def validator(args: typing.Dict):
        return set(expected_args).issubset(args)

    return validator
