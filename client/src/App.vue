<template>
  <div id="app">
    <top-nav></top-nav>
  </div>
</template>

<script>
// @ is an alias to /src
import TopNav from "@/components/TopNavigation.vue";

export default {
  name: "app",
  components: {
    TopNav
  }
};
</script>

<style lang="scss">
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
