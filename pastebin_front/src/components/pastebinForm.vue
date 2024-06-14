<script setup>
import { ref } from "vue";
import axios from 'axios';
import router from "../routes/index";

let selectedFile = null;
const alias = ref("");
const baseUrl = "http://localhost:4000";

function handleFileChange(event) {
  selectedFile = event.target.files[0];
}

async function createPaste() {
  const formData = new FormData();
  formData.append('Alias', alias.value);
  formData.append('Contain', selectedFile);
  
  try {
    await axios.post(`${baseUrl}/bins`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    router.push({ name: 'pastDetails', params: { id: alias.value } });
  } catch(err) {
    console.log(err);
  }
}
</script>

<template>
  <div class="pastebin-form">
    <p>Create new paste</p>
    <div class="pastebin-file">
      <input type="file" @change="handleFileChange" ref="fileInput" />
    </div>
    <div class="pastebin-alias">
      <label for="alias-input">Alias:</label>
      <input id="alias-input" v-model="alias" type="text" />
    </div>

    <button @click="createPaste" class="create-paste-button">
      Create paste
    </button>
    <router-view/>
  </div>
</template>

<style>
.pastebin-form {
  gap: 20px;
  width: 35%;
  border: 4px solid aliceblue;
  margin: auto;
  display: grid;
  padding: 12px;
  border-radius: 8px;
  justify-items: center;
  justify-content: center;

  .pastebin-file {
    cursor: pointer;
  }

  .pastebin-alias {
    gap: 10px;
    display: flex;
  }

  .create-paste-button {
    width: 112px;
    border: 1px solid aliceblue;
    cursor: pointer;
    line-height: 24px;
    border-radius: 8px;
    background-color: aliceblue;
  }
}
</style>
