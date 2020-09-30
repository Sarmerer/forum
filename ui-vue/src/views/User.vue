<template>
  <div class="info-col">
    <div class="card">
      <h3 class="primary">
        {{ user.display_name }} <b-badge v-if="user.role == 2" variant="primary">Admin</b-badge>
      </h3>
      <p>User info</p>
    </div>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
export default {
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
};
</script>
