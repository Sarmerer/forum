#!/bin/bash

git add .
git commit -m "deploy"
git subtree push --prefix . geroku master