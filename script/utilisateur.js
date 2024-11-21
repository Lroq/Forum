document.addEventListener("DOMContentLoaded", function() {
    var switchTabBtn = document.getElementById('switchTabBtn');
    var editProfileBtn = document.getElementById('editProfileBtn');
    function showTab(tabId) {
        var tabs = document.querySelectorAll('.tab');
        tabs.forEach(function(tab) {
            if (tab.id === tabId) {
                tab.classList.add('active');
                tab.style.display = 'block';
            } else {
                tab.classList.remove('active');
                tab.style.display = 'none';
            }
        });
    }
    if (switchTabBtn) {
        switchTabBtn.addEventListener('click', function() {
            var currentTab = document.querySelector('.tab.active');
            var currentTabId = currentTab ? currentTab.id : null;
            if (currentTabId === "activityTab") {
                showTab("profileTab");
            } else if (currentTabId === "profileTab") {
                showTab("activityTab");
            }
        });
    }
    if (editProfileBtn) {
        editProfileBtn.addEventListener('click', function() {
            var editProfileForm = document.getElementById('editProfileForm');
            if (editProfileForm) {
                editProfileForm.style.display = (editProfileForm.style.display === 'block') ? 'none' : 'block';
            }
        });
    }
    showTab("profileTab");
});