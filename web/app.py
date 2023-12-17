from flask import Flask, render_template, make_response, send_from_directory, request
import redis
import hashlib

app = Flask(__name__)
r = redis.Redis(host='localhost', port=6379, db=0)


@app.route('/')
def index():
    return render_template('index.html')

@app.route('/service-worker.js')
def sw():
    response=make_response(
            send_from_directory('static', 'js/service-worker.js'))
    response.headers['Content-Type'] = 'application/javascript'
    return response

@app.route('/manifest.json')
def manifest():
    response=make_response(
        send_from_directory('static', 'manifest.json'))
    return response

@app.route('/app-icon.png')
def icon():
    response=make_response(send_from_directory('static', 'icon-192x192.png'))
    return response



@app.route('/register', methods=['POST'])
def submit():
    r.setex(request.form['targetURL'] + "_$$_" + request.form['subscription'], 172800, request.form['keywords'])
    return "모니터링 등록 완료"


if __name__ == '__main__':
        app.run(host='0.0.0.0', port=2080)

