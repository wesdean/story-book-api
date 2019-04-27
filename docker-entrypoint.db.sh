#!/usr/bin/env bash

psql -U postgres -c "drop database if exists storybook;"
psql -U postgres -c "create database storybook;"
psql -U postgres storybook -c "create user storybook with password 'storybook';"
#psql -U postgres storybook -c "grant all privileges on all tables in schema public to storybook;"

for i in /migrations/*_up.sql; do
    [[ -f "$i" ]] || break
    psql -U storybook storybook -f "$i"
done