### Develop
```shell

mkdir <new lambda python handler func name> && cd {name}
# notice: 
# if u want dev in cur env dir, u could do this, but need to rm .venv when deploy 
# or dev in global env dir
python3 -m venv .venv
source .venv/bin/activate
pip install boto3
pip install pyjwt
pip freeze > requirements.txt
```

### Install
```shell
git clone github.com/weedge/craftsman 
cd craftsman/cloud/aws/cdk/serverless-openai-chatbot/server/lambda/python/login
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```


