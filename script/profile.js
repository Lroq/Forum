function fetchUserPosts() {
    fetch('/user_posts')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(posts => {
            var userPostsDiv = document.getElementById('userPosts');
            if (userPostsDiv) {
                userPostsDiv.innerHTML = '';
                if (posts === null || posts.length === 0) {
                    userPostsDiv.innerHTML = '<p>No posts currently published !</p>';
                    return;
                }
                posts.forEach(post => {
                    var postDiv = document.createElement('div');
                    postDiv.classList.add('post');
                    postDiv.innerHTML = `
                        <h3>${post.Title}</h3>
                        <p>${post.Content}</p>
                        ${post.Image ? `<img src="/uploads/${post.Image}" alt="Image">` : ''}
                        ${post.Gif ? `<img src="/uploads/${post.Gif}" alt="GIF">` : ''}
                        <p><strong>----</strong></p>
                    `;
                    userPostsDiv.appendChild(postDiv);
                });
            }
        })
        .catch(error => {
            console.error('Error fetching user posts:', error);
            alert('An error occurred while fetching user posts.');
        });
}
function fetchUserComments() {
    fetch('/user_comments')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(comments => {
            var userCommentsDiv = document.getElementById('userComments');
            if (userCommentsDiv) {
                userCommentsDiv.innerHTML = '';
                if (comments === null || comments.length === 0) {
                    userCommentsDiv.innerHTML = '<p>No comments currently published !</p>';
                    return;
                }
                comments.forEach(comment => {
                    var commentDiv = document.createElement('div');
                    commentDiv.classList.add('comment');
                    commentDiv.innerHTML = `
                        <h3>Commentaire:${comment.Content}</h3>
                        <p>Post: ${comment.PostTitle}</p>
                        <p><strong>----</strong></p>
                    `;
                    userCommentsDiv.appendChild(commentDiv);
                });
            }
        })
        .catch(error => {
            console.error('Error fetching user comments:', error);
            alert('An error occurred while fetching user comments.');
        });
}
function fetchUserLikes() {
    fetch('/user_likes')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(likes => {
            var userLikesDiv = document.getElementById('userLikes');
            if (userLikesDiv) {
                userLikesDiv.innerHTML = '';
                if (likes === null || likes.length === 0) {
                    userLikesDiv.innerHTML = '<p>No likes currently registered !</p>';
                    return;
                }
                likes.forEach(like => {
                    var likeDiv = document.createElement('div');
                    likeDiv.classList.add('userLikes');
                    likeDiv.innerHTML = `
                        <p>Post: ${like.PostTitle}</p>
                    `;
                    userLikesDiv.appendChild(likeDiv);
                });
            }
        })
        .catch(error => {
            console.error('Error fetching user likes:', error);
            alert('An error occurred while fetching user likes.');
        });
}
function fetchUserDislikes() {
    fetch('/user_dislikes')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(dislikes => {
            var userDislikesDiv = document.getElementById('userDislikes');
            if (userDislikesDiv) {
                userDislikesDiv.innerHTML = '';
                if (dislikes === null || dislikes.length === 0) {
                    userDislikesDiv.innerHTML = '<p>No dislikes currently registered !</p>';
                    return;
                }
                dislikes.forEach(dislike => {
                    var dislikeDiv = document.createElement('div');
                    dislikeDiv.classList.add('userDislikes');
                    dislikeDiv.innerHTML = `
                        <p>Post: ${dislike.PostTitle}</p>
                    `;
                    userDislikesDiv.appendChild(dislikeDiv);
                });
            }
        })
        .catch(error => {
            console.error('Error fetching user dislikes:', error);
            alert('An error occurred while fetching user dislikes.');
        });
}
function showTab(tabId) {
    var tabs = document.querySelectorAll('.tab');
    tabs.forEach(function(tab) {
        tab.classList.remove('active');
        tab.style.display = 'none';
    });
    var selectedTab = document.getElementById(tabId);
    if (selectedTab) {
        selectedTab.classList.add('active');
        selectedTab.style.display = 'block';
    }
}
function deleteComment() {
    const commentId = document.getElementById('commentId').value;
    
    fetch(`/delete_comment/${commentId}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.text();
    })
    .then(data => {
        alert(data); // Display success message or handle as needed
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to delete comment');
    });
}
// Call the functions to fetch user data and initialize the first tab
showTab("profileTab");
fetchUserPosts();
fetchUserComments();
fetchUserLikes();
fetchUserDislikes();