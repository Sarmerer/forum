<template>
  <div class="info-col">
    <div class="card">
      <img
        style="width: 100px; height:100px; border-radius: 200px"
        :src="user.avatar"
        alt="avatar"
        srcset=""
      />
      <h3 class="primary">
        {{ user.display_name }} <b-badge v-if="user.role == 2" variant="primary">Admin</b-badge>
      </h3>
      <p>User info</p>
    </div>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data() {
    return {
      user: {},
    };
  },
  mounted() {
    this.getUser();
  },
  updated() {
    this.getUser();
  },
  methods: {
    async getUser() {
      return await axios
        .get("user", { params: { ID: this.$route.params.id } })
        .then((response) => {
          this.user = response.data.data;
        })
        .catch((error) => {
          console.log(error);
        });
    },
  },
};
</script>
