document.addEventListener('DOMContentLoaded', () => {
  const websocket = new WebSocket('ws://localhost:4000/api/v1/ws');
  const chat = document.querySelector('ol.chat');
  const input = document.querySelector('input#type');

  const addMessage = ({ from, text, avatar }) => {
    const li = document.createElement('li');
    li.classList.add(from);
    const div = document.createElement('div');
    div.classList.add('avatar');
    const img = document.createElement('img');
    img.src = avatar;
    const card = document.createElement('div');
    card.classList.add('msg');
    const content = document.createElement('p');
    content.innerHTML = text;
    const time = document.createElement('time');
    time.innerText = new Date().toISOString().substr(11, 5);

    card.appendChild(content);
    card.appendChild(time);

    div.appendChild(img);
    li.appendChild(div);
    li.appendChild(card);
    chat.appendChild(li);

    window.scrollBy(0, 300);
  };

  const createMessage = (from, text) => ({
    avatar: `/img/${from}.png`,
    context: {},
    from,
    text,
  });

  websocket.onopen = () => {
    websocket.onmessage = (event) => {
      const response = JSON.parse(event.data);
      const assistantMessage = createMessage(
        'assistant',
        response.text || response.output.generic[0].text
      );
      addMessage(assistantMessage);
    };

    input.addEventListener('keypress', (event) => {
      if (event.which === 13 && event.target.value?.length) {
        const { value: text } = event.target;
        const userMessage = createMessage('user', text);
        addMessage(userMessage);
        websocket.send(
          JSON.stringify({
            text: userMessage.text,
          })
        );
        input.value = '';
      }
    });
  };
});
