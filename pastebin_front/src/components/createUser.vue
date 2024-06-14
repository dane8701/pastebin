<script setup>

let email = ref('');
let password = ref('');
let error = ref(false);

function signIn() {
    if (email === '' || password === '') {
        error = true;
        return;
    }
    // Appel de la route
}

function logIn() {
    if (email === '' || password === '') {
        error = true;
        return;
    };
    // Appel de la route 
}

</script>

<template>
    <div class="auth-form">
        <div class="email-input">
            <label for="email">Email:</label>
            <input id="email" v-model="email" type="text"/>
        </div>
        <div class="password-input">
            <label for="password">Password:</label>
            <input type="password" id="password" required/>
        </div>
        <div class="submit-button">
            <button @click="signIn">Sign In</button>
            <button @click="logIn">Log In</button>
        </div>
        <p v-if="error">Username or password is wrong !</p>
    </div>
</template>

<style>

</style>