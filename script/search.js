document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.querySelector('input[name="search"]');
    const posts = document.querySelectorAll('#latestPosts td');

    searchInput.addEventListener('input', function() {
        const searchText = searchInput.value.toLowerCase();

        posts.forEach(post => {
            const title = post.querySelector('h2 a').textContent.toLowerCase();
            if (title.includes(searchText)) {
                post.style.display = '';
            } else {
                post.style.display = 'none';
            }
        });
    });
});