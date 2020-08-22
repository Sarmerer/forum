import { Component } from '@angular/core';

@Component({
	selector: 'app-root',
	templateUrl: './app.component.html',
	styles: [ `.navbar-right { margin-right: 0px !important}` ]
})
export class AppComponent {
	title = 'forum';

	constructor() {}
}
