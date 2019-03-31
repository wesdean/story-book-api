#!/usr/bin/env bash

psql -c "create database storybook;"
psql storybook -c "create user storybook with password 'storybook';"
psql storybook -c "grant select on all tables in schema public to storybook;"