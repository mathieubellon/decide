<script setup>
import { ref, onMounted } from 'vue';

const user = ref({});
const fetchData = async () => {
  try {
    const response = await fetch('/me');
    const data = await response.json();
    console.log(data);
    user.value = data;
  } catch (error) {
    console.error(error);
  }
};
onMounted(() => {
  fetchData();
});
</script>

<template>
  <main>
    <h1>Hello, World 👋!</h1>
    <img :src="user.avatarURL" alt="">
    <p>This is a basic HTML file.</p>
    <a href="/login/google">Auth with Google</a>
    <a href="/login/github">Auth with Github</a>
    <a href="/logout">Logout</a><a href="/ideas">Ideas</a>
    <button @click="fetchData">Fetch API</button>
    <router-link to="/about">About</router-link>
    <pre>{{ user }}</pre>
  </main>
</template>
