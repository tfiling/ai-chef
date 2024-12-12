import tkinter as tk
from tkinter import messagebox, scrolledtext

import yaml


def wait_for_user(msg="Click OK to continue") -> None:
    root = tk.Tk()
    root.withdraw()
    try:
        messagebox.showinfo("Confirmation", msg)
    finally:
        root.destroy()

def ask_yes_no(message: str):
    root = tk.Tk()
    root.withdraw()
    try:
        return messagebox.askyesno("Question", message)
    finally:
        root.destroy()


def get_user_input(message: str, default="") -> str:
    root = tk.Tk()
    root.title(message)

    text = scrolledtext.ScrolledText(root)
    text.pack(padx=10, pady=10)

    result = [""]

    def on_ok():
        result[0] = text.get('1.0', 'end-1c')
        root.destroy()

    tk.Button(root, text="OK", command=on_ok).pack(pady=(0, 10))

    # Wait for user's input
    root.mainloop()
    return result[0] or default

def get_yaml_from_user(message: str):
    try:
        tasks_yaml = get_user_input(message)
        return yaml.safe_load(tasks_yaml)
    except yaml.YAMLError as ex:
        raise yaml.YAMLError(f"Error parsing YAML file: {ex}")
