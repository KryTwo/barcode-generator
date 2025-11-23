#!/bin/bash

echo "Начало скрипта"

# Включаем режим трассировки
set -x

# Команды, которые будут выводиться в консоль
/home/kry/go/bin/fyne-cross windows -arch=amd64 --app-id bcgen.myapp
cp -f /home/kry/GolandProjects/barcode/fyne-cross/bin/windows-amd64/barcode.exe /home/kry/sharing

# Выключаем режим трассировки, если нужно
set +x

echo "Конец скрипта"
