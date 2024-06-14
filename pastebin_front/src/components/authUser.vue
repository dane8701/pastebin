<script>
import { ref } from "vue";
import router from "../routes/index";

export default {
  name: "authUser",
  setup() {
    const email = ref('');
    const password = ref('');
    const error = ref(false);

    function signIn() {
        console.log(this.$router);
      if (email.value === '' || password.value === '') {
        error.value = true;
        return;
      }
      router.push({ name: "newPast" });
      // Appel de la route
    }

    function logIn() {
        console.log(this.$router);
      if (email.value === '' || password.value === '') {
        error.value = true;
        return;
      }
      router.push({ name: "newPast" });
      // Appel de la route 
    }

    return {
      email,
      password,
      error,
      signIn,
      logIn
    };
  }
};
</script>

<template>
    <div class="auth-form">
        <div class="email-input">
            <label for="email">Email:</label>
            <input id="email" v-model="email" type="text"/>
        </div>
        <div class="password-input">
            <label for="password">Password:</label>
            <input type="password" v-model="password" id="password" required/>
        </div>
        <div class="submit-button">
            <button @click="signIn">Sign In</button>
            <button @click="logIn">Log In</button>
        </div>
        <p v-if="error">Username or password is wrong !</p>
    </div>
    <router-view/>
</template>

<style>

</style>