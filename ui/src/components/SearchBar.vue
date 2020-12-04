<template>
  <b-input-group>
    <b-input-group-prepend>
      <b-button @click="$emit('toggle-search')" variant="dark">
        <b-icon-arrow-left></b-icon-arrow-left>
      </b-button>
    </b-input-group-prepend>
    <b-form-input></b-form-input>
    <b-input-group-append>
      <b-button variant="outline-info" v-b-popover.top="'Popover!'">
        <b-icon-info></b-icon-info>
      </b-button>
    </b-input-group-append>
  </b-input-group>
</template>
