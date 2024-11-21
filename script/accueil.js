document.addEventListener('DOMContentLoaded', () => {
    const newPostButton = document.getElementById('newPostButton');

    newPostButton.addEventListener('click', (event) => {
        document.getElementById('postForm').classList.remove('hidden');
    });

    const closeButton = document.getElementById('closeButton');
    closeButton.addEventListener('click', () => {
        document.getElementById('postForm').classList.add('hidden');
    });

    fetch('/check-authorization')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            if (data.role === 'GUEST') {
                disableLikeDislikeButtons();
            }
        })
        .catch(error => {
            console.error('Authorization check failed:', error);
        });

    function disableLikeDislikeButtons() {
        const likeButtons = document.querySelectorAll('form[action="/like"] button[type="submit"]');
        likeButtons.forEach(button => {
            button.disabled = true;
        });

        const dislikeButtons = document.querySelectorAll('form[action="/dislike"] button[type="submit"]');
        dislikeButtons.forEach(button => {
            button.disabled = true;
        });
    }
});
