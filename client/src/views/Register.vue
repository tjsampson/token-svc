<template>
  <div>
    <RegisterUser></RegisterUser>
  </div>
</template>

<script>
// @ is an alias to /src
import RegisterUser from "@/components/RegisterUser.vue";

export default {
  name: "login",
  components: {
    RegisterUser
  }
};
</script>
