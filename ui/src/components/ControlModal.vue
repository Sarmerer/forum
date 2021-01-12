<template>
  <div>
    <b-modal
      :id="modalID"
      centered
      body-bg-variant="darker"
      header-bg-variant="darker"
      hide-header
      hide-footer
      no-fade
    >
      <b-row class="mb-1" @click="$bvModal.hide(modalID)">
        <b-col @click="$emit('edit-event')"><span>Edit</span></b-col>
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
      body-bg-variant="darker"
      hide-header
      hide-footer
      no-fade
      title="Are you sure?"
    >
      <p>Are you sure?</p>
      <b-row align-h="end">
        <b-col cols="end" class="mr-3">
          <b-button
            size="sm"
            class="mr-2"
            variant="outline-danger"
            @click="$emit('delete-event'), $bvModal.hide(modalID + 'confirm')"
          >
            Yes
          </b-button>
          <b-button
            size="sm"
            variant="success"
            @click="$bvModal.hide(modalID + 'confirm')"
          >
            Forget it
          </b-button>
        </b-col>
      </b-row>
    </b-modal>
  </div>
</template>
<script>
export default {
  name: "control-modal",
  props: {
    modalID: { type: String, required: true },
  },
};
</script>
