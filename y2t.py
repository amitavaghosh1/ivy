#!/usr/bin/env python3
from yaml import safe_load
from enum import Enum
import argparse

def render_yaml(param_name, env, values):
    headers = values[0].keys()
    print("| Param Key | Value | Type | Managed By |")
    print("|-----------|-------|------|------------|")

    for row in values:
        typ = "SecureString"
        if row["secured"] == "false" or row["secured"] == False: typ = "String"
        

        print("|/{param_name}/{env}{row_name} | | {typ} | Manual".format(param_name=param_name, env=env, row_name=row["name"], typ=typ))


def build_tfvars(param_name, env, values):
    obj = { "name": env, "parameters": [] }
    for row in values:
        typ = "SecureString"
        if row["secured"] == "false" or row["secured"] == False: typ = "String"

        obj["parameters"].append({ "name": row["name"], "type": typ, "value": row["value"]})

    return obj

def render_tfvars(obj, render_secure):
    strs = ["{",
            "  name = \"{}\"".format(obj["name"]),
            "  parameters = ["]

    for row in obj["parameters"]:
        name = row["name"]
        if name.startswith("/"): name = name.removeprefix("/")

        if row["type"] == "SecureString" and not render_secure: continue

        tfobj = ["  {",
                 "    name = \"{}\"".format(name),
                 "    type = \"{}\"".format(row["type"]),
                 "    value = \"{}\"".format(row["value"]),
                 "  },"]
        strs.append("\n".join(tfobj))

    strs.extend(["]", "},"])
    
    print("\n".join(strs))


def format_to_yaml(data, param_name, render_secure):
    for env in data[param_name]:
        render_yaml(param_name, env, data[param_name][env])

def format_to_tfvars(data, param_name, render_secure):
    for env in data[param_name]:
        render_tfvars(build_tfvars(param_name, env, data[param_name][env]), render_secure)


class OutputTypes(Enum):
    table = 'table'
    tfvars = 'tfvars'

    def __str__(self):
        return self.value

# Instantiate the parser
parser = argparse.ArgumentParser(description='convert config yaml to different format')
parser.add_argument("--file", "-f", type=str, help="yaml file to parse")
parser.add_argument("--format", "-t", type=OutputTypes, choices=list(OutputTypes), help="specify formats. table, tfvar")
parser.add_argument("--secured", "-s", action='store_true', help="render secured string params")
args = parser.parse_args()

with open(args.file) as f:
    x = safe_load(f)

param_name = next(iter(x))

if args.format == OutputTypes.tfvars:
    format_to_tfvars(x, param_name, args.secured)
else:
    format_to_yaml(x, param_name, args.secured)


