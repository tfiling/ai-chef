import typing
from dataclasses import dataclass, field

@dataclass
class PredefinedPrompt:
    name: str
    prompt_template: str
    args_validator: typing.Callable[[typing.Dict[str, typing.Any]],bool] = None
    # TODO - consider deprecating additional_user_contents_reader - using file based user input
    additional_user_contents_reader: typing.Callable[[], typing.Dict[str, typing.Any]] = lambda: {}
    prefix: typing.List[str] = field(default_factory=list)
    postfix: typing.List[str] = field(default_factory=list)

def format_prompt(predefined_prompt: PredefinedPrompt, **kwargs):
    print(f"formatting prompt {predefined_prompt.name}")
    kwargs.update(predefined_prompt.additional_user_contents_reader())
    if predefined_prompt.args_validator and not predefined_prompt.args_validator(kwargs):
        raise ValueError(f"{predefined_prompt.name} received invalid args: {kwargs}")
    return (" ".join(predefined_prompt.prefix) + "\n" +
            predefined_prompt.prompt_template.format(**kwargs) + "\n" +
            " ".join(predefined_prompt.postfix))
