#!/bin/bash
echo enter secret
read secret
curl http://89.108.99.231:5001/webhook/sudo-siteforreg-fork-release/$secret