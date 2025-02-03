const container = document.querySelector('.container');
const registerBtn = document.querySelector('.register-btn');
const loginBtn = document.querySelector('.login-btn');

registerBtn.addEventListener('click', () => {
    container.classList.add('active');
})

loginBtn.addEventListener('click', () => {
    container.classList.remove('active');
})

const toggleLoginPasswordIcon = document.getElementById('toggle-login-password-icon');
const loginPasswordInput = document.getElementById('login-password');

toggleLoginPasswordIcon.addEventListener('click', function() {
    // Toggle password visibility
    if (loginPasswordInput.type === 'password') {
        loginPasswordInput.type = 'text'; // Show password
        toggleLoginPasswordIcon.classList.remove('bxs-show'); // Remove eye icon
        toggleLoginPasswordIcon.classList.add('bxs-hide');
    } else {
        loginPasswordInput.type = 'password'; // Hide password
        toggleLoginPasswordIcon.classList.remove('bxs-hide'); // Remove eye icon
        toggleLoginPasswordIcon.classList.add('bxs-show');
    }
});

const toggleRegisterPasswordIcon = document.getElementById('toggle-register-password-icon');
const loginRegisterInput = document.getElementById('register-password');

toggleRegisterPasswordIcon.addEventListener('click', function() {
    // Toggle password visibility
    if (loginRegisterInput.type === 'password') {
        loginRegisterInput.type = 'text'; // Show password
        toggleRegisterPasswordIcon.classList.remove('bxs-show'); // Remove eye icon
        toggleRegisterPasswordIcon.classList.add('bxs-hide');
    } else {
        loginRegisterInput.type = 'password'; // Hide password
        toggleRegisterPasswordIcon.classList.remove('bxs-hide'); // Remove eye icon
        toggleRegisterPasswordIcon.classList.add('bxs-show');
    }
});