<template>
    <div class="login-view">
        <h1 class="text-lg font-black text-red-400">Login</h1>
        <img src="../assets/github.svg" alt="">
        <div class="social-buttons">
            <a class="loginButton" href="/login/google">
                Login with Google
            </a>
            <a class="loginButton" href="/login/github">
                <img src="../assets/github.svg" alt="">
                Login with GitHub
            </a>
        </div>
    </div>
</template>

<script setup>

function loginWithGoogle() {
    console.log('Login with Google');
}
function loginWithGitHub() {
    console.log('Login with GitHub');
}

</script>

<style scoped>
.login-view {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
}

.social-buttons {
    margin-top: 20px;
}

a.loginButton {
    margin: 10px;
    padding: 15px 25px;
    background-color: #f5f5f5;
    border-color: aquamarine;
    border-radius: 5px;
    cursor: pointer;
    color: black;
}
</style>