import $ from 'jquery';
import React from "react";
import {
    AggregateStatusCount,
    Card,
    CardTitle,
    CardGrid,
    CardBody,
    Icon,
    Row,
    Col,
    AggregateStatusNotifications,
    AggregateStatusNotification,
    Masthead,
    MastheadCollapse,
    MastheadDropdown,
    MenuItem,
    ListGroup,
    ListView,
    ListViewItem, ListViewInfoItem,
} from "patternfly-react";

import "./App.css"

import HorizontalNav from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNav";
import HorizontalNavMenu from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNavMenu";
import HorizontalNavMenuItem from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNavMenuItem";

class App extends React.Component {
    render() {
        return <div>
            <HorizontalNav>
                <Masthead
                    title="IoT simulator console"
                    navToggle={false}
                    thin
                />
                <HorizontalNavMenu>
                    <HorizontalNavMenuItem id="overview" title="Overview"/>
                </HorizontalNavMenu>
            </HorizontalNav>
            <div className="container-fluid">
                <Home/>
            </div>
        </div>
    }
}

export default App;

class Home extends React.Component {

    constructor(props) {
        super(props);
        this.api = "//" + window.location.host + "/api";
        this.state = {
            overview: {
                tenants: []
            },
        };
        this.refreshData();
    }

    refreshData() {
        fetch(this.api + "/overview")
            .then(result => {
                return result.json()
            })
            .then(data => {
                console.log(data);
                this.setState({overview: data});
            })
    }

    componentDidMount() {
        this.interval = setInterval(() => this.refreshData(), 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    renderConsumers(tenant) {
        if (tenant.consumers == null) {
            return
        }
        return tenant.consumers.map(function (consumer, i) {
            return (
                <ListViewItem
                    heading="consumer"
                    leftContent={
                        <Icon
                            type="pf"
                            name={(consumer.messagesPerSecond != null && consumer.messagesPerSecond > 0) ? "ok" : "error-circle-o"}/>
                    }
                    description={consumer.type}
                    additionalInfo={[
                        <ListViewInfoItem>
                            <Icon type="fa" name="cubes"/>
                            <strong>{consumer.replicas} Pods</strong>
                        </ListViewInfoItem>,
                        <ListViewInfoItem>
                            <Icon type="fa" name="angle-double-down"/>
                            <strong title="msgs/s" data-toggle="tooltip" data-placement="top">
                                {(consumer.messagesPerSecond != null) ? consumer.messagesPerSecond.toFixed(0) : "␀"} received
                            </strong>
                        </ListViewInfoItem>
                    ]}
                >
                </ListViewItem>
            )
        })
    }

    renderProducers(tenant) {
        if (tenant.producers == null) {
            return
        }
        return tenant.producers.map(function (producer, i) {
            return (
                <ListViewItem
                    heading="producer"
                    leftContent={
                        <Icon
                            type="pf"
                            name={(producer.messagesPerSecondFailed != null && producer.messagesPerSecondFailed <= 0) ? "ok" : "error-circle-o"}/>
                    }
                    description={producer.type + " / " + producer.protocol}
                    additionalInfo={[
                        <ListViewInfoItem>
                            <Icon type="fa" name="cubes"/>
                            <strong>{producer.replicas} Pods</strong>
                        </ListViewInfoItem>,
                        <ListViewInfoItem>
                            <Icon type="fa" name="angle-double-up"/>
                            <strong>
                                <span title="msgs/s" data-toggle="tooltip" data-placement="top">
                                    {(producer.messagesPerSecondSent != null) ? producer.messagesPerSecondSent.toFixed(0) : "␀"} sent</span>
                            </strong>
                        </ListViewInfoItem>,
                        <ListViewInfoItem>
                            <Icon type="fa" name="clock-o"/>
                            <strong>
                                <span
                                    title="msgs/s configured" data-toggle="tooltip"
                                    data-placement="top">{producer.messagesPerSecondConfigured}</span>&nbsp;→&nbsp;
                                <span
                                    title="msgs/s scheduled" data-toggle="tooltip"
                                    data-placement="top">{(producer.messagesPerSecondScheduled != null) ? producer.messagesPerSecondScheduled.toFixed(0) : "␀"}</span>
                            </strong>
                        </ListViewInfoItem>,
                        <ListViewInfoItem>
                            <Icon type="fa" name="exclamation-triangle"/>
                            <strong>
                                <span
                                    title="msgs/s failed" data-toggle="tooltip"
                                    data-placement="top">{(producer.messagesPerSecondFailed != null) ? producer.messagesPerSecondFailed.toFixed(0) : "␀"}</span>&nbsp;/&nbsp;
                                <span
                                    title="msgs/s errored" data-toggle="tooltip"
                                    data-placement="top">{(producer.messagesPerSecondErrored != null) ? producer.messagesPerSecondErrored.toFixed(0) : "␀"}</span>
                            </strong>
                        </ListViewInfoItem>
                    ]}
                >
                </ListViewItem>
            );
        })
    }

    render() {
        const o = this;
        return (
            <div>
                {
                    this.state.overview.tenants.map(function (tenant, i) {
                        return (
                            <ListGroup>
                                <h2>{tenant.name}</h2>
                                <ListView>
                                    {o.renderConsumers(tenant)}
                                    {o.renderProducers(tenant)}
                                </ListView>
                            </ListGroup>
                        )
                    })
                }
            </div>
        );
    }
}
