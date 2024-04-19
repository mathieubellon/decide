<template>
  <div class="about">
    <h1>This is an about page</h1>
    <a href="/">Go home</a>
    <router-link to="/">Go home</router-link>
  </div>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
