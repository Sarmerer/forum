<template>
  <div>
    <b-table
      class="table"
      :items="users"
      :fields="fields"
      head-variant="dark"
      small
      :busy="fetching"
      table-variant="dark"
    >
      <template #table-busy>
        <div class="text-center text-info my-2">
          <b-spinner class="align-middle"></b-spinner>
          <strong> Fetching...</strong>
        </div>
      </template>
      <template #cell(avatar)="data">
        <b-img-lazy
          width="30px"
          height="30px"
          :src="data.item.avatar"
          rounded
          alt=""
        ></b-img-lazy>
      </template>
    </b-table>
  </div>
</template>
<script>
import api from "@/api/api";

export default {
  data() {
    return {
      users: [],
      fields: [
        { key: "id", label: "ID", sortable: true },
        { key: "avatar", label: "Avatar" },
        { key: "username", label: "username", sortable: true },
        { key: "alias", label: "Name", sortable: true },
        { key: "email", label: "Email", sortable: true },
        { key: "role", label: "Role", sortable: true },
        { key: "created", label: "Created", sortable: true },
      ],
      fetching: true,
    };
  },
  created() {
    this.getUsers();
  },
  methods: {
    async getUsers() {
      return await api
        .get("/users")
        .then((response) => {
          this.users = response.data.data || [];
          let u = this.users[0];
          this.users[0] = {
            id: u.id,
            avatar: u.avatar,
            username: u.username,
            alias: u.alias,
            email: u.email,
            role: u.role,
            created: u.created,
          };
          this.fetching = false;
        })
        .catch(console.log);
    },
  },
};
</script>
<style lang="scss" scoped>
.table {
  margin-top: 5%;
  width: 70%;
  margin-left: 15%;
  margin-right: 15%;
}
</style>
