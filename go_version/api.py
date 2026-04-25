from flask import Flask
app = Flask(__name__)

@app.route("/hello")
def hello():
    return "Bonjour !"

@app.route("/info")
def info():
    return "Informations sur l'API"

if __name__ == "__main__":
    app.run()