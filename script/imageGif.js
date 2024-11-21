function previewImage(event) {
    var input = event.target;
    var preview = document.getElementById('imagePreview');
    
    preview.innerHTML = '';
    if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function(e) {
            var img = document.createElement('img');
            img.src = e.target.result;
            img.alt = 'Preview Image';
            img.style.maxWidth = '200px';
            preview.appendChild(img);
        };
        reader.readAsDataURL(input.files[0]);
    }
}

function previewGif(event) {
    var input = event.target;
    var preview = document.getElementById('gifPreview');
    
    preview.innerHTML = '';
    if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function(e) {
            var img = document.createElement('img');
            img.src = e.target.result;
            img.alt = 'Preview GIF';
            img.style.maxWidth = '200px';
            preview.appendChild(img);
        };
        reader.readAsDataURL(input.files[0]);
    }
}