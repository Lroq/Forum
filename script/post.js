document.addEventListener('DOMContentLoaded', function() {
    const famousForums = document.getElementById('famousForums');
    console.log('famousForums element:', famousForums);

    fetch('http://localhost:8080/top_posts')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(titles => {
            console.log('Fetched titles:', titles);
            if (titles && titles.length > 0) {
                titles.forEach(title => {
                    const item = document.createElement('li');
                    item.textContent = title;
                    console.log('Appending item:', item);
                    famousForums.appendChild(item);
                });
            } else {
                const item = document.createElement('li');
                item.textContent = 'No famous posts found.';
                famousForums.appendChild(item);
            }
        })
        .catch(error => console.error('Error fetching top posts:', error));
});