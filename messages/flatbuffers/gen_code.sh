# Python
echo "flatc: Python"
flatc --python -o ../clients/py-client/ ./messages.fbs

# Go
echo "flatc: Go"
flatc --go -o ../wtf-main/ ./messages.fbs

# C sharp
# flatc --csharp ../clients/unity-client/ ./messages.fbs

