<template>
  <div class="parent d-flex">
    <p>
      Following accounts have identical fields. Would you like to merge them
      into one account?
    </p>
    <div class="flex-container">
      <span
        v-for="(account, index) in accounts"
        :key="index"
        :class="`wrapper flex-${isMobile() ? 'column' : 'row'}`"
      >
        <div class="account-card bg-darker">
          <b-avatar
            :src="account.avatar"
            size="50px"
            :badge="account.oauth_provider != ''"
            badge-variant="dark"
          >
            <template v-if="account.oauth_provider != ''" #badge>
              <b-icon :icon="account.oauth_provider"></b-icon>
            </template>
          </b-avatar>
          <h3 class="text-white-50">{{ account.username }}</h3>
          <small v-if="account.last_active" class="text-white-50">
            Active: <time-ago :datetime="account.last_active"></time-ago>
          </small>
        </div>
        <b-icon
          v-if="index < accounts.length - 1"
          :icon="`arrow-${isMobile() ? 'down' : 'right'}`"
          class="rounded-circle bg-darker p-4"
          scale="1.5"
          shift-v="8"
          shift-h="-8"
        ></b-icon>
      </span>
    </div>
    <div>
      <b-input-group class="mt-3" size="sm">
        <b-input-group-prepend>
          <b-button
            v-b-popover.hover.left="
              'Password forw an existing account, to prove it\'s you'
            "
            variant="outline-dark"
          >
            <b-icon-info></b-icon-info>
          </b-button>
        </b-input-group-prepend>
        <b-form-input placeholder="Password" type="password" v-model="password">
        </b-form-input>
        <b-input-group-append>
          <b-button
            :disabled="!password.length"
            variant="outline-light"
            @click="submitMerge"
          >
            Merge
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </div>
  </div>
</template>
<script>
import TimeAgo from "../components/TimeAgo.vue";
import { mapActions } from "vuex";

export default {
  components: { TimeAgo },
  props: {
    accounts: Array,
  },
  data() {
    return {
      password: "",
    };
  },
  created() {
    if (this.accounts?.length !== 2) return this.$router.push("/");
  },
  methods: {
    ...mapActions({
      merge: "auth/merge",
    }),
    submitMerge() {
      this.merge({
        password: this.password,
        merger: this.accounts[0],
        merged: this.accounts[1],
      }).then(() => {
        let error = this.authError?.data || this.authError?.data;
        if (error) {
          this.$bvToast.toast(error.message, {
            title: "Oops!",
            variant: "danger",
            solid: true,
          });
          this.$store.commit("auth/setAuthError", null);
        }
      });
    },
  },
};
</script>
<style lang="scss" scoped>
.flex-container {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
  margin-top: 20px;
  gap: 10px;
}

.parent {
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  padding-left: 10%;
  padding-right: 10%;
  margin-top: 20px;
}

.wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
}

.account-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 20px;
  width: 200px;
  height: 150px;
}
</style>
