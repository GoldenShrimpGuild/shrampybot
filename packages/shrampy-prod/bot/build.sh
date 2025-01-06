#!/bin/bash

set -e

python -m venv virtualenv
source virtualenv/bin/activate
python -m pip install --upgrade pip
pip install -U -r requirements.txt
deactivate
