<script>
import pastebinForm from "./components/pastebinForm.vue";
import router from "./routes/index";

export default {
  name: "App",
  components: {
    pastebinForm,
  },
  router,

  methods: {
    goToLogin() {
      this.$router.push({ name: 'authUser' });
    },

    goToPastList() {
      this.$router.push({ name: 'pasteList'});
    }
  }
};
</script>

<template>
  <div class="app">
    <div class="buttons">
      <button v-show="$route.name === undefined" class="button-auth" @click="goToLogin">Log In/Sign In</button>
      <button v-show="$route.name === undefined" class="button-auth" @click="goToPastList">List of paste</button>
    </div>
    <pastebinForm v-if="$route.name === undefined"/>
    <RouterView />
  </div>
</template>

<style>
#app {
  .app {
    gap: 10px;
    display: grid;
  }

  .home {
    gap: 24px;
    display: grid;
    justify-content: center;
  }

  .buttons {
    display: flex;
    gap: 8px;

    button {
      width: 112px;
      border: 1px solid aliceblue;
      cursor: pointer;
      line-height: 24px;
      border-radius: 8px;
      background-color: aliceblue;
    }
  }
}
</style>
