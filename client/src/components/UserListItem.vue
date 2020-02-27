<template>
  <tr>
    <td style="text-align:left">
      {{ userDetails.id }}
    </td>
    <td style="text-align:left">
      {{ userDetails.uid }}
    </td>
    <td style="text-align:left">{{ userDetails.email }}</td>
  </tr>
</template>

<script>
export default {
  name: "userlistitem",
  props: {
    name: String,
    userDetails: {}
  }
};
</script>
