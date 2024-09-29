import { Component } from '@angular/core';
import { RequestService } from '@uxwb/ngx-services';
import { links } from '../../../environments/links';
import { AdblockListItem } from './models';

@Component({
  selector: 'app-adblock',
  templateUrl: './adblock.component.html',
  styleUrls: ['./adblock.component.scss'],
})
export class AdblockComponent {

  list: AdblockListItem[] = [];

  constructor(private readonly requestService: RequestService) {
  }

  tabChanged(data: string): void {
    switch (data) {
      case 'List Links':
        this.loadList();
        break;
    }
  }

  private loadList(): void {
    this.requestService.get(links.blacklist_adblock_list).subscribe(value => {
      this.list = value as AdblockListItem[];
    });
  }

  showHideList(index: number) {
    const id = this.list[index].id;
    const active = !(this.list[index].deleted_at === null);
    this.requestService.post(links.blacklist_adblock_list_active, { id, active })
      .subscribe({
        next: (v) => console.log(v),
        error: (e) => console.log(e),
        complete: ()=>console.log('a'),
      });
  }
}
