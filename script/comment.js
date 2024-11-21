fetch('/check-authorization')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            if (data.role === 'GUEST') {
                disableCommentButtons();
            }
        })
        .catch(error => {
            console.error('Authorization check failed:', error);
        });

    function disableCommentButtons() {
        const commentButton = document.querySelectorAll('#commentButton');
        commentButton.forEach(button => {
            button.disabled = true;
        });
    }