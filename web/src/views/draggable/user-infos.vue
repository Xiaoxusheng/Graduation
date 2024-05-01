<template>
  <div class="main-container">
    <TableBody ref="tableBody">
      <template #header>
        <TableHeader
            :show-filter="false"
            title=""
            @search="onSearch"
            @reset-search="onResetSearch"
        >
          <template #search-content>
            <a-form :model="{}" layout="inline">
              <a-form-item v-for="item of conditionItems" :key="item.key" :label="item.label">
                <template v-if="item.render">
                  <FormRender :formItem="item" :render="item.render"/>
                </template>
                <template v-else>
                  <template v-if="item.type === 'input'">
                    <a-input v-model="item.value.value" :placeholder="item.placeholder"/>
                  </template>
                </template>
              </a-form-item>
            </a-form>
          </template>
        </TableHeader>
      </template>
      <template #default>
        <a-card :style="{ marginBottom: '20px',height:'520px',  width:'100%'}" hoverable title="员工信息">
          <a-table
              :bordered="{ wrapper: true, cell: true }"
              :columns="tableColumns"
              :data="dataList"
              :loading="tableLoading"
              :pagination="false"
              :row-selection="{ selectedRowKeys, showCheckedAll }"
              :rowKey="rowKey"
              @selection-change="onSelectionChange"
          >
            <template #columns>
              <a-table-column
                  v-for="item of tableColumns"
                  :key="item.key"
                  :align="item.align"
                  :data-index="(item.key as string)"
                  :fixed="item.fixed"
                  :title="(item.title as string)"
                  :width="item.width"
              >
                <template v-if="item.key === 'index'" #cell="{ rowIndex }">
                  {{ rowIndex + 1 }}
                </template>
                <template v-else-if="item.key === 'department_id'" #cell="{ record }">
                  <a-tag :color="record.department_id === 1 ? 'green' : 'blue'">
                    {{ storedMap.get(record.department_id as number) }}
                  </a-tag>
                </template>
                <template v-else-if="item.key === 'avatar'" #cell="{ record }">
                  <a-avatar
                      :autocapitalize="30"
                      :style="{ backgroundColor: 'var(--color-primary-light-1)' }"
                  >
                    {{ record.name }}
                  </a-avatar>
                </template>
                <template v-else-if="item.key === 'sex'" #cell="{ record }">
                  <a-tag v-if="record.sex === 1" color="blue" size="small">男</a-tag>
                  <a-tag v-else color="blue" size="small">女</a-tag>
                </template>
                <template v-else-if="item.key === 'status'" #cell="{ record }">
                  <a-tag color="purple" size="small">{{ record.status == 1 ? '正常' : '离职' }}</a-tag>
                </template>
                <template v-else-if="item.key === 'pass'" #cell="{ record }">
                  <a-tag v-if="record.pass === 1" color="blue" size="small">通过</a-tag>
                  <a-tag v-else-if="record.pass === 2" color="red" size="small">未通过</a-tag>
                  <a-tag v-else color="purple" size="small">未审核</a-tag>
                </template>
                <template v-else-if="item.key === 'url'" #cell="{ record }">
                  <a-image
                      :preview-props="{
                     actionsLayout: ['rotateRight', 'zoomIn', 'zoomOut'],
               }"
                      :src="record.url"
                      width="100"
                  />
                </template>
                <template v-else-if="item.key === 'actions' " #cell="{ record }">
                  <a-space>
                    <a-button v-if="record.is_examine!==1" size="mini" type="primary" @click="onUpdateItem(record)">
                      审核
                    </a-button>
                    <a-button v-else disabled size="mini" type="primary" @click="onUpdateItem(record)">
                      已审核
                    </a-button>
                  </a-space>
                </template>
              </a-table-column>
            </template>
          </a-table>
        </a-card>

      </template>
      <template #footer>
      </template>
    </TableBody>
  </div>
</template>

<script lang="ts">
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn,} from '@/hooks/table'
import FormRender from '@/components/FormRender'
import {FormItem, ModalDialogType} from '@/types/components'
import {Form, Input, Message} from '@arco-design/web-vue'
import {defineComponent, h, onMounted, ref} from 'vue'
import useUserStore from "@/store/modules/user";
import ModalDialog from "@/components/ModalDialog.vue";
import {get} from "@/api/http";
import {employerInfo} from "@/api/url";

const conditionItems: Array<FormItem> = [
  {
    key: '工号',
    label: '员工工号',
    type: 'input',
    placeholder: '请输入员工工号',
    value: ref(''),
    reset: function () {
      this.value.value = ''
    },
    render: (formItem: FormItem) => {
      return h(Input, {
        placeholder: '请输入员工工号',
        modelValue: formItem.value.value,
        'onUpdate:modelValue': (value) => {
          formItem.value.value = value
        },
      })
    },
  },
]
export default defineComponent({
  name: 'UserInfo',
  components: {
    ModalDialog,
    FormRender,
  },
  setup() {
    const searchForm = ref()
    const pagination = usePagination(doRefresh)
    const {selectedRowKeys, onSelectionChange, showCheckedAll} = useRowSelection()
    const table = useTable()
    const rowKey = useRowKey('id')
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '员工姓名',
        key: 'name',
        dataIndex: 'name',
      },
      {
        title: '工号',
        key: 'uid',
        dataIndex: 'uid',
      },
      {
        title: '性别',
        key: 'sex',
        dataIndex: 'sex',
      },
      {
        title: '部门',
        key: 'department_id',
        dataIndex: 'department_id',
      },
      {
        title: '状态',
        key: 'status',
        dataIndex: 'status',
      },
      {
        title: '生日',
        key: 'birthday',
        dataIndex: 'birthday',
      },
      {
        title: '联系电话',
        key: 'phone',
        dataIndex: 'phone',
      },
      {
        title: '职位',
        key: 'position',
        dataIndex: 'position',
      },
      {
        title: 'ip地址',
        key: 'IP',
        dataIndex: 'IP',
      },
      {
        title: '入职时间',
        key: 'CreatedAt',
        dataIndex: 'CreatedAt',
      },

    ])
    const userStore = useUserStore()
    const actionTitle = ref('请假审核')
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const submitLoading = ref(false)
    const storedMapString = localStorage.getItem('departmentMap');
    const storedMapArray = JSON.parse(storedMapString);
    const storedMap = new Map(storedMapArray)


    let formRef = ref<typeof Form>()
    let load = false

    function doRefresh() {
      table.tableLoading.value = false
    }

    function onSearch() {
      if (conditionItems[0].value.value == '') {
        Message.error("值不能为空！")
        return
      }
      get({
        url: employerInfo,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            uid: conditionItems[0].value.value
          }
        },
      }).then(res => {
        const date = new Date(res.data.birthday);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        res.data.birthday = `${year}-${month}-${day} `
        const CreatedAt = new Date(res.data.CreatedAt);
        const year1 = CreatedAt.getFullYear();
        const month1 = String(CreatedAt.getMonth() + 1).padStart(2, "0");
        const day1 = String(CreatedAt.getDate()).padStart(2, "0");
        res.data.CreatedAt = `${year1}-${month1}-${day1} `
        table.handle(res.data)
        pagination.setTotalSize(res.data.length || 10)
        load = true
      }).catch(error => {
        Message.error(error.message)
      })


    }

    function onResetSearch() {
      conditionItems.forEach((it) => {
        it.reset ? it.reset() : (it.value.value = '')
      })
    }


    // 审核
    function onUpdateItem(record: any) {

    }

    //   弹窗
    function onDataFormConfirm() {
      //


    }

    onMounted(doRefresh)
    return {
      ...table,
      rowKey,
      pagination,
      searchForm,
      tableColumns,
      conditionItems,
      onSearch,
      onResetSearch,
      selectedRowKeys,
      showCheckedAll,
      actionTitle,
      modalDialogRef,
      submitLoading,
      formRef,
      storedMap,
      load,
      onSelectionChange,
      onUpdateItem,
      onDataFormConfirm,
    }
  }

  ,
})
</script>

<style lang="less" scoped>
.form-wrapper {
  font-size: small;
}

.avatar-container {
  position: relative;
  width: 30px;
  height: 30px;
  margin: 0 auto;
  vertical-align: middle;

  .avatar {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }

  .avatar-vip {
    border: 2px solid #cece1e;
  }

  .vip {
    position: absolute;
    top: 0;
    right: -9px;
    width: 15px;
    transform: rotate(60deg);
  }
}

.gender-container {
  .gender-icon {
    width: 20px;
  }
}
</style>
