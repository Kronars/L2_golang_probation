from flask import Flask, request

app = Flask(__name__)

@app.route('/', methods=['POST'])
def index():
    message = request.get_data().decode('utf-8')
    print(message)
    return 'Message received: ' + message

if __name__ == '__main__':
    app.run(debug=True)
