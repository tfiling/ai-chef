import tkinter as tk
from tkinter import messagebox

def wait_for_user(msg="Click OK to continue"):
    root = tk.Tk()
    root.withdraw()
    messagebox.showinfo("Confirmation", msg)
    root.destroy()
