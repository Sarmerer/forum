<template>
  <div>
    <b-modal
      :id="modalID"
      centered
      hide-footer
      hide-header
      body-bg-variant="dark"
      no-fade
    >
      <b-row class="mb-1" @click="$bvModal.hide(modalID)">
        <b-col @click="call(editCallback)"><span>Edit</span></b-col>
        <b-col cols="end">
          <b-icon-x class="h4 mr-3"> </b-icon-x>
        </b-col>
      </b-row>
      <b-row
        @click="$bvModal.hide(modalID), $bvModal.show(modalID + 'confirm')"
      >
        <b-col><span>Delete</span></b-col>
      </b-row>
    </b-modal>
    <b-modal
      :id="modalID + 'confirm'"
      centered
      body-bg-variant="dark"
      header-bg-variant="dark"
      footer-bg-variant="dark"
      no-fade
      title="Are you sure?"
    >
      <template #modal-footer="{ ok,cancel}">
        <b-button
          size="sm"
          variant="danger"
          @click="call(deleteCallback), ok()"
        >
          Yes
        </b-button>
        <b-button size="sm" variant="outline-success" @click="cancel()">
          Forget it
        </b-button>
      </template>
    </b-modal>
  </div>
</template>
<script>
export default {
  name: "control-modal",
  props: {
    deleteCallback: { type: Object, required: true },
    editCallback: { type: Object, required: true },
    modalID: { type: String, required: true },
  },
  methods: {
    call(prop) {
      return prop.args ? prop.callback(...prop.args) : prop.callback();
    },
  },
};
</script>
