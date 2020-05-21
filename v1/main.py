import json
import requests
from threading import Thread

hello_list = [  # from https://www.ethnolink.com.au/how-to-say-hello-in-50-different-languages/
    '(Afrikaans) Goeie dag',
    '(Albanian) Tungjatjeta',
    '(Arabic) Ahlan bik',
    '(Bengali) Nomoskar',
    '(Bosnian) Selam',
    '(Burmese) Mingala ba',
    '(Chinese) Nín hao',
    '(Croatian) Zdravo',
    '(Czech) Nazdar',
    '(Danish) Hallo',
    '(Dutch) Hallo',
    '(Filipino) Helo',
    '(Finnish) Hei',
    '(French) Bonjour',
    '(German) Guten Tag',
    '(Greek) Geia!',
    '(Hebrew) Shalóm',
    '(Hindi) Namasté',
    '(Hungarian) Szia',
    '(Indonesian) Hai',
    '(Iñupiaq) Kiana',
    '(Irish) Dia is muire dhuit',
    '(Italian) Buongiorno',
    '(Japanese) Kónnichi wa',
    '(Korean) Annyeonghaseyo',
    '(Lao) Sabai dii',
    '(Latin) Ave',
    '(Latvian) Es mīlu tevi',
    '(Malay) Selamat petang',
    '(Mongolian) sain baina uu',
    '(Nepali) Namaste',
    '(Norwegian) Hallo.',
    '(Persian) Salâm',
    '(Polish) Witajcie',
    '(Portuguese) Olá',
    '(Romanian) Salut',
    '(Russian) Privét',
    '(Samoan) Talofa',
    '(Serbian) ćao',
    '(Slovak) Nazdar',
    '(Slovene) Zdravo',
    '(Spanish) Hola',
    '(Swahili) Jambo',
    '(Swedish) Hej',
    '(Tagalog) Halo',
    '(Thai) Sàwàtdee kráp',
    '(Turkish) Merhaba',
    '(Ukrainian) Pryvít',
    '(Urdu) Adaab arz hai',
    '(Vietnamese) Chào',
]

config = json.loads(open('config.json', 'r').read())


def verify_web_hook(form):
    if not form or form.get('token') != config['SLACK_TOKEN']:
        raise ValueError('Invalid request/credentials')


def handle_request(request, response_url):
    import random  # move import to here to defer loading until needed

    command = ""
    valid_command = False
    # default to returning an error statement

    if (request):
        commandElements = request.split()
        command = commandElements[0]
        params = commandElements[1:]

    if (command == 'hello'):
        valid_command = True
        message = random.choice(hello_list)
    else:
        message = "hello_bot doesn't know what '{cmd}' means".format(
            cmd=command)

    target_url = ''
    if (valid_command == True):
        # Sends the response back to the channel
        target_url = config['WEBHOOK_URL']
        slack_data = {
            'text': message
        }

    else:
        # Sends the response back to the requester only
        target_url = response_url
        slack_data = {
            'response_type': "ephemeral",
            'text': message
        }

    # Send the result of the command back to Slack
    response = requests.post(
        target_url, data=json.dumps(slack_data),
        headers={'Content-Type': 'application/json'}
    )


def hello_bot(request):
    if request.method != 'POST':
        return 'Only POST requests are accepted', 405

    verify_web_hook(request.form)
    response_url = request.form.get('response_url')

    # Create a thread to handle the request and respond immediately
    # A response time of longer than 3 seconds causes a timeout error message in Slack
    thr_response = Thread(target=handle_request,
                          args=[request.form['text'], response_url])
    thr_response.start()

    # return empty string to reduce spam
    return ''
