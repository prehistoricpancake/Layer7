document.addEventListener('DOMContentLoaded', () => {
    const messageArea = document.getElementById('messageArea');
    const messageInput = document.getElementById('messageInput');
    const sendButton = document.getElementById('sendButton');
    const connectionStatus = document.getElementById('connectionStatus');
    let ws = null;

    function connect() {
        ws = new WebSocket('ws://localhost:8080/ws');

        ws.onopen = () => {
            connectionStatus.textContent = 'Connected';
            connectionStatus.classList.add('connected');
        };

        ws.onclose = () => {
            connectionStatus.textContent = 'Disconnected';
            connectionStatus.classList.remove('connected');
            // Try to reconnect after 5 seconds
            setTimeout(connect, 5000);
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            connectionStatus.textContent = 'Error';
            connectionStatus.classList.remove('connected');
        };

        ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            displayMessage(message, 'received');
        };
    }

    function displayMessage(message, type) {
        const messageElement = document.createElement('div');
        messageElement.classList.add('message', type);
        messageElement.textContent = message.content;
        messageArea.appendChild(messageElement);
        messageArea.scrollTop = messageArea.scrollHeight;
    }

    function sendMessage() {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            alert('Not connected to server');
            return;
        }

        const content = messageInput.value.trim();
        if (content) {
            const message = {
                content: content,
                timestamp: new Date().toISOString()
            };

            ws.send(JSON.stringify(message));
            displayMessage(message, 'sent');
            messageInput.value = '';
        }
    }

    sendButton.addEventListener('click', sendMessage);
    messageInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });

    // Initial connection attempt
    connect();
});