<template>
  <div>User Logout</div>
</template>

<script>
export default {
  name: "logout",
  data() {
    return {};
  },
  created() {
    const { dispatch } = this.$store;
    dispatch("user/logout");
  }
};
</script>

<style scoped></style>
