// Initialize global variables for the socket connection and UI elements.
let socket = null
let o = document.getElementById('output')
let userField = document.getElementById('username')
let messageField = document.getElementById('message')

// Define behavior for when the user attempts to leave the page.
window.onbeforeunload = function () {
  // Notify the server that the user has left the chat.
  let jsonData = {}
  jsonData['action'] = 'left'
  socket.send(JSON.stringify(jsonData))
}

// Once the DOM is fully loaded, configure websocket connection and event listeners.
document.addEventListener('DOMContentLoaded', function () {
  // Establish the websocket connection with the server.
  socket = new ReconnectingWebSocket('ws://127.0.0.1:8080/ws', null, {
    debug: true,
    reconnectInterval: 3000
  })

  // Indicators for the connection status.
  const offline = `<span class="badge bg-danger">Not connected</span>`
  const online = `<span class="badge bg-success">Connected</span>`
  let statusDiv = document.getElementById('status')

  // Event handler when the websocket connection is opened.
  socket.onopen = () => {
    console.log('Successfully connected')
    statusDiv.innerHTML = online
  }

  // Event handler when the websocket connection is closed.
  socket.onclose = () => {
    console.log('Connection closed')
    statusDiv.innerHTML = offline
  }

  // Event handler for any websocket errors.
  socket.onerror = (error) => {
    console.log('There was an error with the websocket connection: ', error)
  }

  // Event handler for receiving messages through the websocket.
  socket.onmessage = (msg) => {
    let data = JSON.parse(msg.data)
    console.log('Action is', data.action)

    // Handle different types of actions received via websocket.
    switch (data.action) {
      case 'list_users':
        // Update the list of online users.
        let ul = document.getElementById('online_users')
        while (ul.firstChild) ul.removeChild(ul.firstChild)

        if (data.connected_users.length > 0) {
          data.connected_users.forEach(function (item) {
            let li = document.createElement('li')
            li.appendChild(document.createTextNode(item))
            ul.appendChild(li)
          })
        }
        break

      case 'broadcast':
        // Display the broadcasted message in the chatbox.
        o.innerHTML += data.message + '<br>'
        break
    }
  }

  // When the username is changed, notify the server of the new username.
  userField.addEventListener('change', function () {
    let jsonData = {}
    jsonData['action'] = 'username'
    jsonData['username'] = this.value
    socket.send(JSON.stringify(jsonData))
  })

  // Event handler for sending a message when the Enter key is pressed.
  messageField.addEventListener('keydown', function (event) {
    if (event.code === 'Enter') {
      // Do not proceed if there is no connection or if the username/message fields are empty.
      if (!socket || userField.value === '' || messageField.value === '') {
        errorMessage('Fill out username and message')
      } else {
        sendMessage()
      }
      // Prevent the default action of the Enter key and stop it from propagating.
      event.preventDefault()
      event.stopPropagation()
    }
  })

  // Click event handler for the 'Send Message' button.
  document.getElementById('sendBtn').addEventListener('click', function () {
    // Do not proceed if the username/message fields are empty.
    if (userField.value === '' || messageField.value === '') {
      errorMessage('Fill out username and message')
    } else {
      sendMessage()
    }
  })
})

// Function to send a message to the server via websocket.
function sendMessage() {
  let jsonData = {}
  jsonData['action'] = 'broadcast'
  jsonData['username'] = userField.value
  jsonData['message'] = messageField.value
  socket.send(JSON.stringify(jsonData))
  messageField.value = '' // Clear the message field after sending.
}

// Function to display error messages using the notie notification library.
function errorMessage(msg) {
  notie.alert({
    type: 'error',
    text: msg
  })
}
