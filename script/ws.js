document.addEventListener('DOMContentLoaded', (event) => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = function(event) {
        console.log("Connected to WebSocket");
    };

    socket.onmessage = function(event) {
        console.log("Received message:", event.data);
        let postData = JSON.parse(event.data);
        let latestPosts = document.getElementById('latestPosts');
        let post = document.createElement('div');
        post.innerHTML = `
            <table>
                <tr>
                    <td>
                        <p id="postID" class="hidden">${postData.id}</p>
                        <h3 id="username">${postData.username}</h3>
                        <h2>${postData.title}</h2>
                        <p>${postData.content}</p>
                        <p>Like : <span id="likeCount">0</span></p>
                        <p>Dislike : <span id="dislikeCount">0</span></p>
                        <form action="/like" method="POST" class="likeForm">
                            <input type="hidden" name="postID" value="${postData.id}">
                            <button id="likeButton" type="submit">Like</button>
                        </form>
                        <form action="/dislike" method="POST" class="dislikeForm">
                            <input type="hidden" name="postID" value="${postData.id}">
                            <button id="dislikeButton" type="submit">Dislike</button>
                        </form>
                        <button type="button" class="commentButton">Comment</button>
                    </td>
                </tr>
            </table>
        `;
        latestPosts.appendChild(post);

        post.querySelector('.likeForm').addEventListener('submit', handleLike);
        post.querySelector('.dislikeForm').addEventListener('submit', handleDislike);
    };

    socket.onclose = function(event) {
        if (event.code === 4001) {
            alert('You are not authorized to create / like / dislike / comment posts. You have to be connected to an account to perform these actions.');
        } else {
            console.log("Disconnected from WebSocket, reason:", event.reason);
        }
    };

    socket.onerror = function(error) {
        console.error("WebSocket error:", error);
    };

    document.getElementById('newPostButton').addEventListener('click', (event) => {
        fetch('/check-authorization')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Not authorized');
                }
                return response.json();
            })
            .then(data => {
                console.log('Authorization check result:', data);
                if (data.role === 'GUEST') {
                    alert('You are not authorized to create / like / dislike / comment posts. You have to be connected to an account to perform these actions.');
                    window.location.href = '/login';
                } else {
                    document.getElementById('postForm').classList.remove('hidden');
                }
            })
            .catch(error => {
                console.error('Authorization check failed:', error);
                window.location.href = '/login';
            });
    });
});
