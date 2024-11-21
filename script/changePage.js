document.addEventListener('DOMContentLoaded', function() {
    const pageSelect = document.getElementById('pageSelect');
    pageSelect.addEventListener('change', function() {
        const userRoleElement = document.getElementById('userRole');
        const userRole = userRoleElement.textContent.trim();
        const selectedValue = this.value;
        let url;
        switch(selectedValue) {
            case 'home':
                url = '/';
                break;
            case 'profile':
                url = `/profile/${userRole}`;
                break;
            case 'contact':
                url = '/contact';
                break;
            case 'about':
                url = '/about';
                break;
            default:
                console.log('Unexpected selectedValue:', selectedValue);
        }
        console.log('Generated URL:', url);
        if (url) {
            window.location.href = url;
        }
    });
});
