import {OnInit, Component, EventEmitter, Input, Output, HostListener} from '@angular/core';
import {ContentedService} from './contented_service';

import {Directory} from './directory';
import * as _ from 'lodash';

@Component({
    selector: 'directory-cmp',
    templateUrl: 'directory.ng.html'
})
export class DirectoryCmp implements OnInit {

    @Input() dir: Directory;
    @Input() previewWidth: number;
    @Input() previewHeight: number;
    @Input() currentViewItem: string;
    @Input() maxRendered: number = 3; // Default setting for how many should be visible at any given time

    // @Output clickEvt: EventEmitter<any>;
    public visibleSet: Array<string>;

    constructor() {

    }

    public ngOnInit() {
        console.log("Directory Component loading up");
    }


    public getVisibleSet(currentItem = this.currentViewItem, max: number = this.maxRendered) {
        this.visibleSet = null;
        this.visibleSet = this.dir.getIntervalAround(currentItem, max, 1);
        return this.visibleSet;
    }

    public imgLoaded(evt) {
        let img = evt.target;
        console.log("Img Loaded", img.naturalHeight, img.naturalWidth, img);
    }
}

