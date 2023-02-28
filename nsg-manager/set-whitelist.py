#!/usr/bin/env python3

from ipwhois import IPWhois
import requests
from azure.identity import ClientSecretCredential
from azure.mgmt.network import NetworkManagementClient
from time import sleep
from azure.core.exceptions import HttpResponseError

"""
Simple class for adding the CIDR your public IP is part of and add it into the NSG Security Rule for not manually doing this all the time manually.
A Service Principal needs to be first created within the subscription and give proper permissions to read NSGs and update Security rules.
Params:
@tenant_id: The Tenant ID of the Service Principal
@app_id: The Service Principal App ID
@sub_id: Subscription ID
@rg_name: Resource Group name where the NSG is created
@nsg_name: NSG Name
@rule_name: Security Rule name from the NSG
"""
class NSGManager:
    def __init__(self, tenant_id, app_id, secret, sub_id, rg_name, nsg_name, rule_name):
        self.tenant_id = tenant_id
        self.app_id = app_id
        self.secret = secret
        self.sub_id = sub_id
        self.rg_name = rg_name
        self.nsg_name = nsg_name
        self.rule_name = rule_name

        requests.packages.urllib3.util.connection.HAS_IPV6 = False  # needed to force IPv4

    def get_cidr(self):
        url = "http://icanhazip.com"
        ip = requests.get(url).text.strip()

        res = IPWhois(ip).lookup_rdap()

        cidr = res["asn_cidr"]

        print("[*] Got CIDR:", cidr)

        return cidr

    def update_rule(self, cidr):
        cred_info = (self.tenant_id, self.app_id, self.secret)
        cred = ClientSecretCredential(*cred_info)
        net_client = NetworkManagementClient(cred, self.sub_id)

        rule = net_client.security_rules.get(
            self.rg_name, self.nsg_name, self.rule_name)
        rule.source_address_prefix = cidr  # update security rule

        try:
            poller = net_client.security_rules.begin_create_or_update(
                self.rg_name, self.nsg_name, self.rule_name, security_rule_parameters=rule)
            while not poller.done():
                print("[*] Status:", poller.status())
                sleep(1)
            print("[*] Security Rule has been updated")
        except HttpResponseError as e:
            print(e)


def main():
     nsg_manager = NSGManager(
    tenant_id="xxxx-xxxx",
    app_id="xxxx-xxxx",
    secret="xxxxx",
    sub_id="xxxxxx",
    rg_name="xxxx",
    nsg_name="xxxx",
    rule_name="xxxxx",
)
     cidr = nsg_manager.get_cidr()
     nsg_manager.update_rule(cidr)
     
if __name__ == "__main__": main()

