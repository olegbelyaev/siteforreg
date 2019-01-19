#!/bin/bash
read -s -p "Enter webhook secret:" secret
echo -e "\nsending..."
curl http://89.108.99.231:5001/webhook/sudo-siteforreg-fork-release/$secret